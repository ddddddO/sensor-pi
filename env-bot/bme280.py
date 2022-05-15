#!/usr/bin/env python3
#coding: utf-8

# copied from https://github.com/SWITCHSCIENCE/samplecodes/blob/master/BME280/Python27/bme280_sample.py
from smbus2 import SMBus

class BME280:
	def __init__(self):
		self.bus_number  = 1
		self.bus = SMBus(self.bus_number)
		self.i2c_address = 0x76

		self.digT = []
		self.digP = []
		self.digH = []

		self.t_fine = 0.0

		# FIXME:
		self.result_T = 0.0
		self.result_P = 0.0
		self.result_H = 0.0

		self.setup()

	def writeReg(self, reg_address, data):
		self.bus.write_byte_data(self.i2c_address,reg_address,data)

	def get_calib_param(self):
		calib = []

		for i in range (0x88,0x88+24):
			calib.append(self.bus.read_byte_data(self.i2c_address,i))
		calib.append(self.bus.read_byte_data(self.i2c_address,0xA1))
		for i in range (0xE1,0xE1+7):
			calib.append(self.bus.read_byte_data(self.i2c_address,i))

		self.digT.append((calib[1] << 8) | calib[0])
		self.digT.append((calib[3] << 8) | calib[2])
		self.digT.append((calib[5] << 8) | calib[4])
		self.digP.append((calib[7] << 8) | calib[6])
		self.digP.append((calib[9] << 8) | calib[8])
		self.digP.append((calib[11]<< 8) | calib[10])
		self.digP.append((calib[13]<< 8) | calib[12])
		self.digP.append((calib[15]<< 8) | calib[14])
		self.digP.append((calib[17]<< 8) | calib[16])
		self.digP.append((calib[19]<< 8) | calib[18])
		self.digP.append((calib[21]<< 8) | calib[20])
		self.digP.append((calib[23]<< 8) | calib[22])
		self.digH.append( calib[24] )
		self.digH.append((calib[26]<< 8) | calib[25])
		self.digH.append( calib[27] )
		self.digH.append((calib[28]<< 4) | (0x0F & calib[29]))
		self.digH.append((calib[30]<< 4) | ((calib[29] >> 4) & 0x0F))
		self.digH.append( calib[31] )

		for i in range(1,2):
			if self.digT[i] & 0x8000:
				self.digT[i] = (-self.digT[i] ^ 0xFFFF) + 1

		for i in range(1,8):
			if self.digP[i] & 0x8000:
				self.digP[i] = (-self.digP[i] ^ 0xFFFF) + 1

		for i in range(0,6):
			if self.digH[i] & 0x8000:
				self.digH[i] = (-self.digH[i] ^ 0xFFFF) + 1

	def readData(self):
		data = []
		for i in range (0xF7, 0xF7+8):
			data.append(self.bus.read_byte_data(self.i2c_address,i))
		pres_raw = (data[0] << 12) | (data[1] << 4) | (data[2] >> 4)
		temp_raw = (data[3] << 12) | (data[4] << 4) | (data[5] >> 4)
		hum_raw  = (data[6] << 8)  |  data[7]

		self.compensate_T(temp_raw)
		self.compensate_P(pres_raw)
		self.compensate_H(hum_raw)

	def compensate_P(self, adc_P):
		global  t_fine
		pressure = 0.0

		v1 = (t_fine / 2.0) - 64000.0
		v2 = (((v1 / 4.0) * (v1 / 4.0)) / 2048) * self.digP[5]
		v2 = v2 + ((v1 * self.digP[4]) * 2.0)
		v2 = (v2 / 4.0) + (self.digP[3] * 65536.0)
		v1 = (((self.digP[2] * (((v1 / 4.0) * (v1 / 4.0)) / 8192)) / 8)  + ((self.digP[1] * v1) / 2.0)) / 262144
		v1 = ((32768 + v1) * self.digP[0]) / 32768

		if v1 == 0:
			return 0
		pressure = ((1048576 - adc_P) - (v2 / 4096)) * 3125
		if pressure < 0x80000000:
			pressure = (pressure * 2.0) / v1
		else:
			pressure = (pressure / v1) * 2
		v1 = (self.digP[8] * (((pressure / 8.0) * (pressure / 8.0)) / 8192.0)) / 4096
		v2 = ((pressure / 4.0) * self.digP[7]) / 8192.0
		pressure = pressure + ((v1 + v2 + self.digP[6]) / 16.0)
		pressure = pressure/100

		self.result_P = float("{:7.2f}".format(pressure))

		print("pressure : %7.2f hPa" % (pressure))

	def compensate_T(self, adc_T):
		global t_fine
		v1 = (adc_T / 16384.0 - self.digT[0] / 1024.0) * self.digT[1]
		v2 = (adc_T / 131072.0 - self.digT[0] / 8192.0) * (adc_T / 131072.0 - self.digT[0] / 8192.0) * self.digT[2]
		t_fine = v1 + v2
		temperature = t_fine / 5120.0

		self.result_T = float("{:-6.2f}".format(temperature))

		result = "temp : {:-6.2f} ℃".format(temperature)
		print(result)

	def compensate_H(self, adc_H):
		global t_fine
		var_h = t_fine - 76800.0
		if var_h != 0:
			var_h = (adc_H - (self.digH[3] * 64.0 + self.digH[4]/16384.0 * var_h)) * (self.digH[1] / 65536.0 * (1.0 + self.digH[5] / 67108864.0 * var_h * (1.0 + self.digH[2] / 67108864.0 * var_h)))
		else:
			return 0
		var_h = var_h * (1.0 - self.digH[0] * var_h / 524288.0)
		if var_h > 100.0:
			var_h = 100.0
		elif var_h < 0.0:
			var_h = 0.0

		self.result_H = float("{:6.2f}".format(var_h))

		print("hum : %6.2f ％" % (var_h))


	def setup(self):
		osrs_t = 1			#Temperature oversampling x 1
		osrs_p = 1			#Pressure oversampling x 1
		osrs_h = 1			#Humidity oversampling x 1
		mode   = 3			#Normal mode
		t_sb   = 5			#Tstandby 1000ms
		filter = 0			#Filter off
		spi3w_en = 0			#3-wire SPI Disable

		ctrl_meas_reg = (osrs_t << 5) | (osrs_p << 2) | mode
		config_reg    = (t_sb << 5) | (filter << 2) | spi3w_en
		ctrl_hum_reg  = osrs_h

		self.writeReg(0xF2,ctrl_hum_reg)
		self.writeReg(0xF4,ctrl_meas_reg)
		self.writeReg(0xF5,config_reg)

	def result(self) -> (float, float, float):
		# print(self.result_T) # e.g. 27.9
		# print(self.result_P) # e.g. 997.31
		# print(self.result_H) # e.g. 57.1
		return self.result_T, self.result_P, self.result_H

if __name__ == '__main__':
	try:
		bme280 = BME280()
		bme280.get_calib_param()
		bme280.readData()
	except KeyboardInterrupt:
		pass
