from kafka import KafkaConsumer
import json
import base64

# consumer = KafkaConsumer('test', group_id= 'group2', bootstrap_servers= ['localhost:9092'], api_version=(2,8,1), value_deserializer=lambda m: json.loads(m.decode('utf-8')))
consumer = KafkaConsumer('dataOut', group_id= 'group2', bootstrap_servers= ['localhost:9092'], api_version=(2,8,1))
for msg in consumer:
    # img = base64.b64decode(msg)
    # obj = json.loads(msg)
    # print(obj.value)
    print(msg.key)
    print(msg.value)
    # img = base64.b64decode(msg.value)
    # file = open('test.jpg','wb')
    # file.write(img)
    # file.close()
    