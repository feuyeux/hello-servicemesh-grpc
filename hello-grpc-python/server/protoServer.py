# encoding: utf-8

import time
import uuid

import grpc
from concurrent import futures

from landing_pb2 import landing_pb2
from landing_pb2 import landing_pb2_grpc

hellos = ["Hello", "Bonjour", "Hola", "こんにちは", "Ciao", "안녕하세요"]


def build_result(data):
    result = landing_pb2.TalkResult()
    result.id = int((time.time()))
    result.type = landing_pb2.OK
    result.kv["id"] = str(uuid.uuid1())
    result.kv["idx"] = data
    index = int(data)
    result.kv["data"] = hellos[index]
    result.kv["meta"] = "PYTHON"
    return result


class LandingServiceServer(landing_pb2_grpc.LandingServiceServicer):
    def talk(self, request, context):
        response = landing_pb2.TalkResponse()
        response.status = 200
        result = build_result(request.data)
        response.results.append(result)
        return response

    def talkOneAnswerMore(self, request, context):
        datas = request.data.split(",")
        for data in datas:
            response = landing_pb2.TalkResponse()
            response.status = 200
            result = build_result(data)
            response.results.append(result)
            yield response

    def talkMoreAnswerOne(self, request_iterator, context):
        response = landing_pb2.TalkResponse()
        response.status = 200
        for request in request_iterator:
            response.results.append(build_result(request.data))
        return response

    def talkBidirectional(self, request_iterator, context):
        for request in request_iterator:
            response = landing_pb2.TalkResponse()
            response.status = 200
            result = build_result(request.data)
            response.results.append(result)
            yield response


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    landing_pb2_grpc.add_LandingServiceServicer_to_server(LandingServiceServer(), server)
    server.add_insecure_port('[::]:9996')
    server.start()
    try:
        while True:
            time.sleep(60 * 60 * 24)
    except KeyboardInterrupt:
        server.stop(0)


if __name__ == '__main__':
    serve()
