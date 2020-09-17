#!/usr/bin/env bash
SCRIPT_PATH="$(
  cd "$(dirname "$0")" >/dev/null 2>&1
  pwd -P
)/"
cd "$SCRIPT_PATH" || exit

USER_CONFIG=~/shop_config/ack_bj
MESH_CONFIG=~/shop_config/asm_bj
alias k="kubectl --kubeconfig $USER_CONFIG"
alias m="kubectl --kubeconfig $MESH_CONFIG"

client_java_pod=$(k get pod -l app=grpc-client-java -n grpc-best -o jsonpath={.items..metadata.name})
client_go_pod=$(k get pod -l app=grpc-client-go -n grpc-best -o jsonpath={.items..metadata.name})
client_node_pod=$(k get pod -l app=grpc-client-node -n grpc-best -o jsonpath={.items..metadata.name})
client_python_pod=$(k get pod -l app=grpc-client-python -n grpc-best -o jsonpath={.items..metadata.name})
for ((i = 1; i <= 4; i++)); do
  k exec "$client_java_pod" -c grpc-client-java -n grpc-best -- java -jar /grpc-client.jar
  echo
  echo
done
for ((i = 1; i <= 4; i++)); do
  k exec "$client_go_pod" -c grpc-client-go -n grpc-best -- ./grpc-client
  echo
  echo
done
for ((i = 1; i <= 4; i++)); do
  k exec "$client_node_pod" -c grpc-client-node -n grpc-best -- node proto_client.js
  echo
  echo
done
for ((i = 1; i <= 4; i++)); do
  k exec "$client_python_pod" -c grpc-client-python -n grpc-best -- sh /grpc-client/start_client.sh
  echo
  echo
done