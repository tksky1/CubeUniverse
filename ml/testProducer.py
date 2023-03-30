from kafka import KafkaProducer
import json
import base64

# 连接kafka
producer = KafkaProducer(bootstrap_servers=['localhost:9092'], api_version=(0,11,5), value_serializer=lambda m: json.dumps(m).encode('utf-8'))
# producer = KafkaProducer(bootstrap_servers=['localhost:9092'])

# 处理图片，转为ascii码
# image = open('cropped_panda.jpg', 'rb')
image = open('pics/code.jpeg', 'rb')
res = base64.b64encode(image.read())
image_str=res.decode('ascii')

# 发送
future = producer.send('picIn' , key= b'112233.png', value= image_str, partition= 0)
result = future.get(timeout= 10)
print(result)