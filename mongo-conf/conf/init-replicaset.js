var cfg = {
    "_id": "rs0",
    "members": [
        {
            "_id": 0,
            "host": "mongo-server.1:27017",
            "priority": 2
        },
        {
            "_id": 1,
            "host": "mongo-server.2:27017",
            "priority": 0
        },
        {
            "_id": 2,
            "host": "mongo-server.3:27017",
            "priority": 0
        }
    ]
}
rs.initiate(cfg)
rs.conf()