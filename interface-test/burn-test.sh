#!/bin/bash

NUM=0

while [ $NUM -le 0 ]; do
  # curl -v --request POST \
  #   --header "Authorization: beb51d947b0358fe9216a410880cdb20299a12788dbaafa425712f026f762a4c" \
  #   --data '{"tasks":[{"boardname":"TinySim","deviceid":"","runtime":30,"filehash":"42834f9027d659f9900e5e6971b950e5541702d06dcc0347c28fa5929d669dd8","clientid":"","taskindex":1}, {"boardname":"TinySim","deviceid":"","runtime":30,"filehash":"42834f9027d659f9900e5e6971b950e5541702d06dcc0347c28fa5929d669dd8","clientid":"","taskindex":2}], "pid": ""}' \
  #   http://kubernetes.tinylink.cn/linklab/device-control-v2/user-service/api/device/burn

  # curl -v --request POST \
  #   --header "Authorization: 423c0889878290e7fc0369089a7bf0284bc5871b1055b3cb803e4d974a805b13" \
  #   --data '{"tasks":[{"boardname":"Haas100","deviceid":"","runtime":240,"filehash":"d6dbfdee8e4173b07677ca537be5584c44e99d1e77fd03b05e511d50e904c265","clientid":"","taskindex":1}], "pid": ""}' \
  #   http://10.214.149.214:31958/api/device/burn
  
    # curl -v --request POST \
    # --header "Authorization: 5b9bf255425362a6bd5bd5a5488f3034cca8921def8ec5ce6d72fb501998e89d" \
    # --data '{"tasks":[{"boardname":"Haas100","deviceid":"","runtime":240,"filehash":"e2dab532eca7202051b5ef17ff1b23bdd71e45282b6e6a2660801041a314adc6","clientid":"","taskindex":1}], "pid": ""}' \
    # http://192.168.88.20:31958/api/device/burn

    # curl -v --request POST \
    # --header "Authorization: 9ca4f83e28648468ed335953905d8877bb63f7f22c2cd20e52e3799912c04707" \
    # --data '{"tasks":[{"boardname":"ESP32DevKitCArduino","deviceid":"","runtime":240,"filehash":"771b3218979fb8ac6b200576d622f146e233341846dbdfdeb104b1735ddc6bf8","clientid":"","taskindex":1}], "pid": ""}' \
    # http://kubernetes.tinylink.cn/linklab/device-control-v2/user-service/api/device/burn

    # d6dbfdee8e4173b07677ca537be5584c44e99d1e77fd03b05e511d50e904c265
    # e2dab532eca7202051b5ef17ff1b23bdd71e45282b6e6a2660801041a314adc6
    # a9d19d9b3e85bf019b58d5bd6bcd626d3086baa7e4feb2e2fccf35c70be0431a
    # 00695e776a905395d2ec83e52600638a6f4b0bd039f493b668559cb0694254fb
    # 8caaf592180b79eec26ed77b9250106d138436a87cb1be85f0ce12e4213142ae
    # 1768b93f6082b884df93d330e3496ba3f4873bf92bd304bc02c2744719f81823
    # 02c6a680882127b3f36d4f5e2fa7ae7815a2eb64971a67b7a5480728ff7a04b9

     curl -v --request POST \
    --header "Authorization: fb70630a5cda4ec40dd62cd669dce0ccb22115ecf9d95d1ec2f0a98feb6b2bbc" \
    --data '{"tasks":[{"boardname":"Python3Exec","deviceid":"","runtime":10,"filehash":"492dbde09a4728ebd34b84b77e4998f7a88c1c92f42b2f9cc7a354a6fd9f01d5","clientid":"","taskindex":1}], "pid": ""}' \
    http://kubernetes.tinylink.cn/linklab/device-control-v2/user-service/api/device/burn

    # curl -v --request POST \
    # --header "Authorization: 9ca4f83e28648468ed335953905d8877bb63f7f22c2cd20e52e3799912c04707" \
    # --data '{"tasks":[{"boardname":"TinySim","deviceid":"","runtime":240,"filehash":"42834f9027d659f9900e5e6971b950e5541702d06dcc0347c28fa5929d669dd8","clientid":"","taskindex":1}], "pid": ""}' \
    # http://kubernetes.tinylink.cn/linklab/device-control-v2/user-service/api/device/burn


    

    # curl -v --request POST \
    # --header "Authorization: c2f7193b949290ebbc519d3edd7c3b0735595c26cc8256228e58d1046ae6570e" \
    # --data '{"tasks":[{"boardname":"ESP32DevKitC","deviceid":"","runtime":30,"filehash":"6ec0d4238b7164784a62f4b163c712c480b45b2a931ed6c2c6b00e4c66890ca1","clientid":"","taskindex":1}], "pid": ""}' \
    # http://kubernetes.tinylink.cn/linklab/device-control-v2/user-service/api/device/burn

    # curl -v --request POST \
    # --header "Authorization: 25ddd65eebdd5127c820bd37e8bcb9f8c493c0e8f2c6c37109ff972c615514b0" \
    # --data '{"tasks":[{"boardname":"ArduinoMega2560","deviceid":"","runtime":30,"filehash":"0b192f1d33d6d6fe6b13c6e34e0ea5a98e6dd97fcc56f0f289d2f1af43378220","clientid":"","taskindex":1}], "pid": "12099"}' \
    # http://kubernetes.tinylink.cn/linklab/device-control-v2/user-service/api/device/burn

    # curl -v --request POST \
    # --header "Authorization: c4cbdfdd226b48d3453e61459fcb38883fa467b2d9ba2470a3d5d0e8beabffe7" \
    # --data '{"tasks":[{"boardname":"TelosB","deviceid":"","runtime":30,"filehash":"5a4b62d3f696647c5d5f7274b9412b6b921dd68d18b7e01c444283713fcc5a39","clientid":"","taskindex":1}], "pid": ""}' \
    # http://kubernetes.tinylink.cn/linklab/device-control-v2/user-service/api/device/burn

  NUM=$(( $NUM + 1 ))
done


