// linklab compile
db = db.getSiblingDB("linklab")

// linklab compile user
db.createUser({ user: "Compile", pwd: "12", roles: [{ role: "readWrite", db: "linklab" }] })

// linklab create compile collection
db.createCollection("compile", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            properties: {
                _id: {
                    bsonType: "objectId",
                    description: "must be a objectId and is required"
                },
                type: {
                    bsonType: "string",
                    description: "must be a string and is required"
                },
                compileType: {
                    bsonType: "string",
                    description: "must be a string and is required"
                },
                boardType: {
                    bsonType: "string",
                    description: "must be a string and is required"
                },
                fileHash: {
                    bsonType: "string",
                    pattern: "^[A-Fa-f0-9]{40}$",
                    description: "must be a string and is required"
                },
                branch: {
                    bsonType: "string",
                    description: "must be a string and is required"
                },
                fileData: {
                    bsonType: "binData",
                    description: "must be a binData and is required"
                },
                output: {
                    bsonType: "binData",
                    description: "must be a binData and is required"
                },
                message: {
                    bsonType: "string",
                    description: "must be a string and is required"
                },
                status: {
                    bsonType: "string",
                    description: "must be a string and is required"
                }
            }
        }
    }
})

db.compile.createIndex({ compileType: 1, boardType: 1, fileHash: 1 }, { unique: true })