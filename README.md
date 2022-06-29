cd kafka-tools
docker build -t registry.cn-beijing.aliyuncs.com/boer/ktools:0.9.8 .
docker push registry.cn-beijing.aliyuncs.com/boer/ktools:0.9.8

```yaml
kubectl create -n mw-kafka -f - <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: kafka-test
  namespace: mw-kafka
spec:
  containers:
  - name: kafka-test
    image: registry.cn-beijing.aliyuncs.com/boer/ktools:0.9.8
    # Just spin & wait forever
    command: [ "/bin/sh", "-c", "--" ]
    args: [ "while true; do sleep 3000; done;" ]
EOF
```

```bash
kubectl -n mw-kafka exec -it kafka-test -- sh
# with auth
ktop -kafkaURL=kafka-headless:9092 -withAuth -auths=admin:admin -topic=test002 -p=6 -r=2
kcat -L -b kafka-headless:29000
kpro -kafkaURL=kafka-headless:9092 -withAuth -auths=admin:admin -topic=test002 -acks=1
kcom -kafkaURL=kafka-headless:9092 -withAuth -auths=admin:admin -topic=test002 -groupID=cgi
# with out auth
ktop -kafkaURL=kafka-headless:29000 -topic=test003 -p=3 -r=1
kcat -L -b kafka-headless:29000
kpro -kafkaURL=kafka-headless:29000 -topic=test003 -acks=1
kcom -kafkaURL=kafka-headless:29000 -topic=test003 -groupID=cgi
```