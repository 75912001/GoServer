# -*- coding: utf-8 -*-
import top.api

#url = "http://gw.api.taobao.com/router/rest"
#port = 80
appkey = 23273583
secret = "28409ec2fdac3a381fe7546f55493900"
req=top.api.AlibabaAliqinFcSmsNumSendRequest()#url,port)
req.set_app_info(top.appinfo(appkey,secret))

req.sms_type="normal"
req.sms_free_sign_name="注册验证"
req.sms_param="{'code':'123','product':'爱车旅'}"
req.rec_num="17721027200"
req.sms_template_code="SMS_2515091"
try:
    resp= req.getResponse()
    print(resp)
except Exception,e:
    print(e)
    print("######")

    
