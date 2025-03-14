import json
import requests
import pandas as pd
import asyncio
import websockets
import time
import threading
#### setup
ldc_url_http = "http://kubernetes.kaiwutech.cn/linklab/device-control-v2/"
ldc_url_https = "https://kubernetes.kaiwutech.cn/linklab/device-control-v2/"
target_board = "ArduinoUno"
#### message
#### result_path
errordevice = "./data/errordevice.csv"
error_list = []
def addgroup_single(boardtype):
    # data = pd.read_csv(save_path)
    # detection_target = []
    # for index,row in data.iterrows():
    #     if target_board in row['boardtype']:
    #         detection_target.append([row['clientid'],'/dev/'+row['deviceid']])

    ##login
    url_login = ldc_url_http+"login-authentication/user/login"
    param_login = {"id":"UserTest","password":"6b51d431df5d7f141cbececcf79edf3dd861c3b4069f0b11661a3eefacbba918"}
    response = requests.post(url_login, data=json.dumps(param_login))
    if response.status_code == 200:
        # 读取并处理返回的JSON数据
        data = response.json()
        # 在这里可以对返回的数据进行进一步处理或断言验证
        print(data)
    else:
        print('API request failed with status code:', response.status_code)
    token = data['data']['token']

    ##device list
    url_upload = ldc_url_http+"user-service/api/device/list"
    parameters = {"boardname": boardtype}
    # data = {
    #     'parameters': json.dumps(parameters)
    # }
    response = requests.get(url_upload,params=parameters,headers={'Authorization': token})

    if response.status_code == 200:
        # 读取并处理返回的JSON数据
        data = response.json()
        # 在这里可以对返回的数据进行进一步处理或断言验证
    else:
        print('API request failed with status code:', response.status_code)
    device_info = data['data']['devices']
    # for list in device_info:
    #     print(list['deviceid'])
    #     list['deviceid'] = list['deviceid'].replace('/dev/','')
    print(device_info)

    for list in device_info:
        url_create = ldc_url_http+""
        url_add = ldc_url_http+"user-service/api/device/linkgroup"
        parameters = {
            "type": "Single"+list['boardname'],
            "devices": [
                {
                    "clientid": list['clientid'],
                    "deviceid": list['deviceid']
                }
            ]
        }
        # data = {
        #     'parameters': json.dumps(parameters)
        # }
        response = requests.post(url_add, data=json.dumps(parameters), headers={'Authorization': token})
        if response.status_code == 200:
            # 读取并处理返回的JSON数据
            data = response.json()
            print(data)
            print('Device: '+list['deviceid']+' to Client: '+list['clientid']+' Success!')
            # 在这里可以对返回的数据进行进一步处理或断言验证
        else:
            print(response.reason)
            print('Failed addition--Device: '+list['deviceid']+' to Client: '+list['clientid'])

def addgroup_double(boardtype,group):
    # data = pd.read_csv(save_path)
    # detection_target = []
    # for index,row in data.iterrows():
    #     if target_board in row['boardtype']:
    #         detection_target.append([row['clientid'],'/dev/'+row['deviceid']])

    ##login
    url_login = ldc_url_http+"login-authentication/user/login"
    param_login = {"id":"UserTest","password":"6b51d431df5d7f141cbececcf79edf3dd861c3b4069f0b11661a3eefacbba918"}
    response = requests.post(url_login, data=json.dumps(param_login))
    if response.status_code == 200:
        # 读取并处理返回的JSON数据
        data = response.json()
        # 在这里可以对返回的数据进行进一步处理或断言验证
        print(data)
    else:
        print('API request failed with status code:', response.status_code)
    token = data['data']['token']

    ##device list
    url_list = ldc_url_http+"user-service/api/device/list"
    parameters = {"boardname": boardtype}
    # data = {
    #     'parameters': json.dumps(parameters)
    # }
    response = requests.get(url_list,params=parameters,headers={'Authorization': token})

    if response.status_code == 200:
        # 读取并处理返回的JSON数据
        data = response.json()
        # 在这里可以对返回的数据进行进一步处理或断言验证
    else:
        print('API request failed with status code:', response.status_code)
    device_info = data['data']['devices']
    # for list in device_info:
    #     print(list['deviceid'])
    #     list['deviceid'] = list['deviceid'].replace('/dev/','')
    print(device_info)

    for list in device_info:
        for list2 in device_info:
            if list2['deviceid']!=list['deviceid']:
                print(list2['deviceid'],list['deviceid'])
                # url_create = ldc_url_http+""
                url_add = ldc_url_http+"user-service/api/device/linkgroup"
                parameters = {
                    "type": group,
                    "devices": [
                        {
                            "clientid": list['clientid'],
                            "deviceid": list['deviceid']
                        },
                        {
                            "clientid": list2['clientid'],
                            "deviceid": list2['deviceid']
                        }
                    ]
                }

                response = requests.post(url_add, data=json.dumps(parameters), headers={'Authorization': token})
                if response.status_code == 200:
                    # 读取并处理返回的JSON数据
                    data = response.json()
                    print(data)
                    print('Device: ' + list['deviceid'] + list2['deviceid'] + ' to Client: ' + list[
                        'clientid'] + ' Success!')
                    # 在这里可以对返回的数据进行进一步处理或断言验证
                else:
                    print(response.reason)
                    print('Failed addition--Device: ' + list['deviceid'] + ' to Client: ' + list['clientid'])
        # data = {
        #     'parameters': json.dumps(parameters)
        # }

def addgroup_double_cut(boardtype,group):
    # data = pd.read_csv(save_path)
    # detection_target = []
    # for index,row in data.iterrows():
    #     if target_board in row['boardtype']:
    #         detection_target.append([row['clientid'],'/dev/'+row['deviceid']])

    ##login
    url_login = ldc_url_http+"login-authentication/user/login"
    param_login = {"id":"UserTest","password":"6b51d431df5d7f141cbececcf79edf3dd861c3b4069f0b11661a3eefacbba918"}
    response = requests.post(url_login, data=json.dumps(param_login))
    if response.status_code == 200:
        # 读取并处理返回的JSON数据
        data = response.json()
        # 在这里可以对返回的数据进行进一步处理或断言验证
        print(data)
    else:
        print('API request failed with status code:', response.status_code)
    token = data['data']['token']

    ##device list
    url_list = ldc_url_http+"user-service/api/device/list"
    parameters = {"boardname": boardtype}
    # data = {
    #     'parameters': json.dumps(parameters)
    # }
    response = requests.get(url_list,params=parameters,headers={'Authorization': token})

    if response.status_code == 200:
        # 读取并处理返回的JSON数据
        data = response.json()
        # 在这里可以对返回的数据进行进一步处理或断言验证
    else:
        print('API request failed with status code:', response.status_code)
    device_info = data['data']['devices']
    # for list in device_info:
    #     print(list['deviceid'])
    #     list['deviceid'] = list['deviceid'].replace('/dev/','')
    print(device_info)

    for i in range(int(len(device_info)/2)):
        list = device_info[i*2]
        if i*2+1<len(device_info):
            list2 = device_info[i*2+1]
            print(list2['deviceid'],list['deviceid'])
            # url_create = ldc_url_http+""
            url_add = ldc_url_http+"user-service/api/device/linkgroup"
            parameters = {
                "type": group,
                "devices": [
                    {
                        "clientid": list['clientid'],
                        "deviceid": list['deviceid']
                    },
                    {
                        "clientid": list2['clientid'],
                        "deviceid": list2['deviceid']
                    }
                ]
            }
            response = requests.post(url_add, data=json.dumps(parameters), headers={'Authorization': token})
            if response.status_code == 200:
                # 读取并处理返回的JSON数据
                data = response.json()
                print(data)
                print('Device: ' + list['deviceid'] + list2['deviceid'] + ' to Client: ' + list[
                    'clientid'] + ' Success!')
                # 在这里可以对返回的数据进行进一步处理或断言验证
            else:
                print(response.reason)
                print('Failed addition--Device: ' + list['deviceid'] + ' to Client: ' + list['clientid'])
        # data = {
        #     'parameters': json.dumps(parameters)
        # }

def create_group(grouptype,device):
    ##login
    url_login = ldc_url_https + "login-authentication/user/login"
    param_login = {"id": "UserTest", "password": "6b51d431df5d7f141cbececcf79edf3dd861c3b4069f0b11661a3eefacbba918"}
    response = requests.post(url_login, data=json.dumps(param_login))
    if response.status_code == 200:
        # 读取并处理返回的JSON数据
        data = response.json()
        # 在这里可以对返回的数据进行进一步处理或断言验证
        print(data)
    else:
        print('API request failed with status code:', response.status_code)
    token = data['data']['token']

    ##device list
    url_create = ldc_url_https + "user-service/api/device/creategroup"
    parameters = {
        "type": grouptype,
        "boards": device
    }
    # data = {
    #     'parameters': json.dumps(parameters)
    # }
    response = requests.post(url_create, data=json.dumps(parameters), headers={'Authorization': token})

    if response.status_code == 200:
        # 读取并处理返回的JSON数据
        data = response.json()
        print(data)
        # 在这里可以对返回的数据进行进一步处理或断言验证
    else:
        print('API request failed with status code:', response.text)
def delete_group(groupid):
    ##login
    url_login = ldc_url_https + "login-authentication/user/login"
    param_login = {"id": "UserTest", "password": "6b51d431df5d7f141cbececcf79edf3dd861c3b4069f0b11661a3eefacbba918"}
    response = requests.post(url_login, data=json.dumps(param_login))
    if response.status_code == 200:
        # 读取并处理返回的JSON数据
        data = response.json()
        # 在这里可以对返回的数据进行进一步处理或断言验证
        print(data)
    else:
        print('API request failed with status code:', response.status_code)
    token = data['data']['token']

    ##device list
    url_delete = ldc_url_https + "user-service/api/device/unlinkgroup"
    parameters = {
        "id": groupid,
    }
    # data = {
    #     'parameters': json.dumps(parameters)
    # }
    response = requests.post(url_delete, data=json.dumps(parameters), headers={'Authorization': token})

    if response.status_code == 200:
        # 读取并处理返回的JSON数据
        data = response.json()
        print(data)
        # 在这里可以对返回的数据进行进一步处理或断言验证
    else:
        print('API request failed with status code:', response.text)
def delete_group(groupid):
    ##login
    url_login = ldc_url_https + "login-authentication/user/login"
    param_login = {"id": "UserTest", "password": "6b51d431df5d7f141cbececcf79edf3dd861c3b4069f0b11661a3eefacbba918"}
    response = requests.post(url_login, data=json.dumps(param_login))
    if response.status_code == 200:
        # 读取并处理返回的JSON数据
        data = response.json()
        # 在这里可以对返回的数据进行进一步处理或断言验证
        print(data)
    else:
        print('API request failed with status code:', response.status_code)
    token = data['data']['token']

    ##device list
    url_delete = ldc_url_https + "user-service/api/device/unlinkgroup"
    parameters = {
        "id": groupid,
    }
    # data = {
    #     'parameters': json.dumps(parameters)
    # }
    response = requests.post(url_delete, data=json.dumps(parameters), headers={'Authorization': token})

    if response.status_code == 200:
        # 读取并处理返回的JSON数据
        data = response.json()
        print(data)
        # 在这里可以对返回的数据进行进一步处理或断言验证
    else:
        print('API request failed with status code:', response.text)
def delete_group_appointed_devicetype(devicetype,grouplist):
    for list in grouplist:
        if devicetype in list['id']:
            groupid = list['id']
            print(groupid)
            delete_group(groupid)

def get_grouplist():
    ##login
    url_login = ldc_url_https + "login-authentication/user/login"
    param_login = {"id": "UserTest", "password": "6b51d431df5d7f141cbececcf79edf3dd861c3b4069f0b11661a3eefacbba918"}
    response = requests.post(url_login, data=json.dumps(param_login))
    if response.status_code == 200:
        # 读取并处理返回的JSON数据
        data = response.json()
        # 在这里可以对返回的数据进行进一步处理或断言验证
        print(data)
    else:
        print('API request failed with status code:', response.status_code)
    token = data['data']['token']

    ##device list
    url_get = ldc_url_https + "user-service/api/device/listlinkgroup"
    # data = {
    #     'parameters': json.dumps(parameters)
    # }
    response = requests.get(url_get, headers={'Authorization': token})

    if response.status_code == 200:
        # 读取并处理返回的JSON数据
        data = response.json()
        print(data['data']['groups'])
        return data['data']['groups']

        # 在这里可以对返回的数据进行进一步处理或断言验证
    else:
        print('API request failed with status code:', response.text)



# create_group("DoubleSmartVilla",["SmartVilla", "SmartVilla"])
# addgroup_single("SmartVilla")
# addgroup_double("SmartVilla","DoubleSmartVilla")
# addgroup_double_cut("SmartVilla","DoubleSmartVilla")
# for list in get_grouplist():
#     print(list)
# delete_group_appointed_devicetype('DoubleESP32DevKitC-',get_grouplist())