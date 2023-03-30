# AI部分

## 前置要求

- Python ≥ 3.7
- pip ≥ 20.0
- 安装tensorflow
    - pip3 install tensorflow

## 文件结构

```bash
ml
├── classifier.py
├── cropped_panda.jpg   // 用于测试的图片
├── docker-compose.yml  // 测试环境相关
├── inception // 机器学习模型
│   ├── classify_image_graph_def.pb
│   ├── cropped_panda.jpg
│   ├── imagenet_2012_challenge_label_map_proto.pbtxt
│   ├── imagenet_synset_to_human_label_map.txt
│   └── LICENSE
├── picClassifier.py    // 主要程序
├── pics
│   └── * // 用于测试的图片
├── README.md
├── start-kafka.sh      // 测试环境相关
├── testConsumer.py     // 测试用的消费者
├── test.jpg            // picClassifier运行时产生的临时文件
└── testProducer.py     // 测试用的生产者
```

## 使用方法

向 **picIn** 这个 topic 发送 k-v（图片名-图片base64编码后的ascii表示），从 ************dataOut************ 读 k-v（图片名-图片识别结果）

Demo代码在 testProducer.py 和 testConsumer.py

- 数据流向
    
    **testProducer** —picIn→ **picClassifier** —dataOut→ **testConsumer**