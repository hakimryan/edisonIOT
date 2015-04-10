#!usr/bin/env python

import mraa
import time
import json
import pycurl
import datetime

try:
	lm35 = mraa.Aio(0)
	while 1:
		lm35val = lm35.read()
		print lm35val
		time.sleep(1)

		##------------------------ kirim data -------------------------
		localtime = datetime.datetime.today()
		data = {"method":"SensorService.LM35", "params":[{"Time":str(localtime), "Temp": str(lm35val)}], "id": "1"}

		c = pycurl.Curl()
		c.setopt(c.URL, "http://10.10.8.50:10000/rpc")
		c.setopt(c.HTTPHEADER, ["Content-Type: application/json"])
		post = json.dumps(data)
		c.setopt(c.POST, 1)
		c.setopt(c.POSTFIELDS, post)
		c.perform()
		##-------------------------------------------------------------

except:
	print "exception!"
