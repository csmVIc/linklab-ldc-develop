// use admin;
db = db.getSiblingDB("admin")
db.createUser({ user: "LinkLab", pwd: "12", roles: [{ role: "root", db: "admin" }] })

// create linklab db
db = db.getSiblingDB("linklab")
db.createUser({ user: "DeviceControl", pwd: "12", roles: [{ role: "readWrite", db: "linklab" }] })
db.createUser({ user: "Client", pwd: "12", roles: [{ role: "read", db: "linklab" }] })
db.createUser({ user: "ProblemSystem", pwd: "12", roles: [{ role: "read", db: "linklab" }] })

db.createCollection("users", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            properties: {
                _id: {
                    bsonType: "objectId",
                    description: "must be a objectId and is required"
                },
                userId: {
                    bsonType: "string",
                    pattern: "^[a-zA-Z][a-zA-Z0-9_-]{3,15}$",
                    description: "must be a string and is required"
                },
                email: {
                    bsonType: "string",
                    pattern: "^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$",
                    description: "must be a string and is required"
                },
                hash: {
                    bsonType: "string",
                    pattern: "^[A-Fa-f0-9]{64}$",
                    description: "must be a string and is required"
                },
                salt: {
                    bsonType: "string",
                    pattern: "^[A-Fa-f0-9]{64}$",
                    description: "must be a string and is required"
                },
                tenantId: {
                    bsonType: "int",
                    description: "must be a int and is required"
                }
            }
        }
    }
})
db.users.createIndex({ userId: 1 }, { unique: true })
db.users.createIndex({ email: 1 }, { unique: true })

db.createCollection("clients")
db.clients.createIndex({ username: 1 }, { unique: true })
// db.clients.insert({ "username": "DeviceControl", "password": "4cbc7eb5c02f16dc12a492cefe5e3fdbe65edb0b3910463c0d5101402e7ca230", "is_superuser": true, "salt": "1e37ade6945020dd65ae01652ea48a8ceb62d7cb2c61f90149dc29ad931febbf", "tenantId": NumberInt(1) })
// db.clients.insert({ "username": "ClientTest", "password": "4cbc7eb5c02f16dc12a492cefe5e3fdbe65edb0b3910463c0d5101402e7ca230", "is_superuser": false, "salt": "1e37ade6945020dd65ae01652ea48a8ceb62d7cb2c61f90149dc29ad931febbf", "tenantId": NumberInt(1) })
db.clients.insert({ "username": "DeviceControl", "password": "4cbc7eb5c02f16dc12a492cefe5e3fdbe65edb0b3910463c0d5101402e7ca230", "is_superuser": true, "salt": "1e37ade6945020dd65ae01652ea48a8ceb62d7cb2c61f90149dc29ad931febbf", "tenantId": { "1": true } })
db.clients.insert({ "username": "ClientTest", "password": "4cbc7eb5c02f16dc12a492cefe5e3fdbe65edb0b3910463c0d5101402e7ca230", "is_superuser": false, "salt": "1e37ade6945020dd65ae01652ea48a8ceb62d7cb2c61f90149dc29ad931febbf", "tenantId": { "1": true } })

db.createCollection("files", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            properties: {
                _id: {
                    bsonType: "objectId",
                    description: "must be a objectId and is required"
                },
                boardName: {
                    bsonType: "string",
                    description: "must be a string and is required"
                },
                fileHash: {
                    bsonType: "string",
                    pattern: "^[A-Fa-f0-9]{64}$",
                    description: "must be a string and is required"
                },
                fileData: {
                    bsonType: "binData",
                    description: "must be a binData and is required"
                }
            }
        }
    }
})
db.files.createIndex({ boardName: 1, fileHash: 1 }, { unique: true })

db.createCollection("boards", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            properties: {
                _id: {
                    bsonType: "objectId",
                    description: "must be a objectId and is required"
                },
                boardName: {
                    bsonType: "string",
                    description: "must be a string and is required"
                },
                boardType: {
                    bsonType: "string",
                    description: "must be a string and is required"
                }
            }
        }
    }
})
db.boards.createIndex({ boardName: 1 }, { unique: true })

db.boards.insert({ "boardName": "ESP32DevKitC", "boardType": "IoTNode" })
db.boards.insert({ "boardName": "ESP32DevKitCArduino", "boardType": "IoTNode" })
db.boards.insert({ "boardName": "ArduinoMega2560", "boardType": "IoTNode" })
db.boards.insert({ "boardName": "TelosB", "boardType": "IoTNode" })
db.boards.insert({ "boardName": "ArduinoUno", "boardType": "IoTNode" })
db.boards.insert({ "boardName": "DeveloperKit", "boardType": "IoTNode" })
db.boards.insert({ "boardName": "Haas100", "boardType": "IoTNode" })
db.boards.insert({ "boardName": "Haas100Python", "boardType": "IoTNode" })
db.boards.insert({ "boardName": "STM32F103C8", "boardType": "IoTNode" })
db.boards.insert({ "boardName": "ArduinoMega2560WithHC06", "boardType": "IoTNode" })
db.boards.insert({ "boardName": "SmartVilla", "boardType": "IoTNode" })
db.boards.insert({ "boardName": "HaasEDUK1", "boardType": "IoTNode" })
db.boards.insert({ "boardName": "nRF52840", "boardType": "IoTNode" })
db.boards.insert({ "boardName": "VirtualDevice", "boardType": "VirtualNode" })
db.boards.insert({ "boardName": "TinySim", "boardType": "LinuxHostNode" })
db.boards.insert({ "boardName": "Python3Exec", "boardType": "LinuxHostNode" })
db.boards.insert({ "boardName": "AIBox", "boardType": "LinuxHostNode" })
// db.boards.insert({ "boardName": "RaspberryPi4B", "boardType": "EdgeNode" })

db.createCollection("boardtypes", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            properties: {
                _id: {
                    bsonType: "objectId",
                    description: "must be a objectId and is required"
                },
                boardType: {
                    bsonType: "string",
                    description: "must be a string and is required"
                },
                allowCmd: {
                    bsonType: "bool",
                    description: "must be a bool and is required"
                }
            }
        }
    }
})
db.boardtypes.createIndex({ boardType: 1 }, { unique: true })

db.boardtypes.insert({ "boardType": "IoTNode", "allowCmd": true })
db.boardtypes.insert({ "boardType": "VirtualNode", "allowCmd": false })
db.boardtypes.insert({ "boardType": "LinuxHostNode", "allowCmd": false })
// db.boardtypes.insert({ "boardType": "EdgeNode", "allowCmd": false })

// linklab create devicelog collection
db.createCollection("devicelog", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            properties: {
                _id: {
                    bsonType: "objectId",
                    description: "must be a objectId and is required"
                },
                userId: {
                    bsonType: "string",
                    description: "must be a string and is required"
                },
                clientId: {
                    bsonType: "string",
                    description: "must be a string and is required"
                },
                devPort: {
                    bsonType: "string",
                    description: "must be a string and is required"
                },
                waitingId: {
                    bsonType: "string",
                    description: "must be a string and is required"
                },
                groupId: {
                    bsonType: "string",
                    description: "must be a string and is required"
                },
                pid: {
                    bsonType: "string",
                    description: "must be a string and is required"
                },
                logs: {
                    bsonType: "array",
                    minItems: 0,
                    items: {
                        bsonType: "string"
                    },
                    description: "must be a array of strings and is required"
                },
                isEnd: {
                    bsonType: "bool",
                    description: "must be a bool and is required"
                },
                startDate: {
                    bsonType: "date",
                    description: "must be a date and is required"
                },
                endDate: {
                    bsonType: "date",
                    description: "must be a date and is required"
                }
            }
        }
    }
})

// linklab create devicelog index
db.devicelog.createIndex({ pid: 1 })
db.devicelog.createIndex({ waitingId: 1 }, { unique: true })

// linklab create grouplog collection
db.createCollection("grouplog", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            properties: {
                _id: {
                    bsonType: "objectId",
                    description: "must be a objectId and is required"
                },
                groupId: {
                    bsonType: "string",
                    description: "must be a string and is required"
                },
                userId: {
                    bsonType: "string",
                    description: "must be a string and is required"
                },
                waitingIds: {
                    bsonType: "array",
                    minItems: 1,
                    items: {
                        bsonType: "string",
                    },
                    description: "must be a array of strings and is required"
                },
                pid: {
                    bsonType: "string",
                    description: "must be a string and is required"
                },
                logs: {
                    bsonType: "array",
                    minItems: 0,
                    items: {
                        bsonType: "string"
                    },
                    description: "must be a array of strings and is required"
                }
            }
        }
    }
})

// linklab create grouplog index
db.grouplog.createIndex({ pid: 1 })
db.grouplog.createIndex({ groupId: 1 }, { unique: true })

// linklab create tenants collection
db.createCollection("tenants", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            properties: {
                _id: {
                    bsonType: "objectId",
                    description: "must be a objectId and is required"
                },
                tenantId: {
                    bsonType: "int",
                    description: "must be a int and is required"
                },
                tenantName: {
                    bsonType: "string",
                    pattern: "^[a-zA-Z][a-zA-Z0-9_-]{3,15}$",
                    description: "must be a string and is required"
                },
                isSystemTenant: {
                    bsonType: "bool",
                    description: "must be a bool and is required"
                }
            }
        }
    }
})

// linklab create tenants index
db.tenants.createIndex({ tenantId: 1 }, { unique: true })

db.tenants.insert({ "tenantId": NumberInt(1), "tenantName": "Linklab", "isSystemTenant": true })
db.tenants.insert({ "tenantId": NumberInt(10000), "tenantName": "Emlab", "isSystemTenant": false })


db.createCollection("podyaml", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            properties: {
                _id: {
                    bsonType: "objectId",
                    description: "must be a objectId and is required"
                },
                fileHash: {
                    bsonType: "string",
                    pattern: "^[A-Fa-f0-9]{64}$",
                    description: "must be a string and is required"
                },
                fileData: {
                    bsonType: "binData",
                    description: "must be a binData and is required"
                }
            }
        }
    }
})

// linklab create podyaml index
db.podyaml.createIndex({ fileHash: 1 }, { unique: true })

db.createCollection("imagebuild", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            properties: {
                _id: {
                    bsonType: "objectId",
                    description: "must be a objectId and is required"
                },
                fileHash: {
                    bsonType: "string",
                    pattern: "^[A-Fa-f0-9]{64}$",
                    description: "must be a string and is required"
                },
                fileData: {
                    bsonType: "binData",
                    description: "must be a binData and is required"
                }
            }
        }
    }
})

// linklab create imagebuild index
db.imagebuild.createIndex({ fileHash: 1 }, { unique: true })

db.createCollection("boardgroups")

// linklab create boardgroups index
db.boardgroups.createIndex({ type: 1 }, { unique: true })


