# encoding: utf-8
import logging
import time
import uuid

import grpc
from concurrent import futures

from landing_pb2 import landing_pb2
from landing_pb2 import landing_pb2_grpc

logger = logging.getLogger('grpc-server')
logger.setLevel(logging.INFO)
console = logging.StreamHandler()
console.setLevel(logging.INFO)
formatter = logging.Formatter('%(asctime)s [%(levelname)s] - %(message)s')
console.setFormatter(formatter)
logger.addHandler(console)

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


def read_headers(method_name, context):
    metadata = context.invocation_metadata()
    metadata_dict = {}
    for c in metadata:
        logger.info("%s HEADER %s:%s", method_name, c.key, c.value)
        metadata_dict[c.key] = c.value
    return metadata_dict


class LandingServiceServer(landing_pb2_grpc.LandingServiceServicer):
    def talk(self, request, context):
        read_headers("TALK", context)
        logger.info("TALK REQUEST: data=%s,meta=%s", request.data, request.meta)
        response = landing_pb2.TalkResponse()
        response.status = 200
        result = build_result(request.data)
        response.results.append(result)
        return response

    def talkOneAnswerMore(self, request, context):
        read_headers("TalkOneAnswerMore", context)
        logger.info("TalkOneAnswerMore REQUEST: data=%s,meta=%s", request.data, request.meta)
        datas = request.data.split(",")
        for data in datas:
            response = landing_pb2.TalkResponse()
            response.status = 200
            result = build_result(data)
            response.results.append(result)
            yield response

    def talkMoreAnswerOne(self, request_iterator, context):
        read_headers("TalkMoreAnswerOne", context)
        response = landing_pb2.TalkResponse()
        response.status = 200
        for request in request_iterator:
            logger.info("TalkMoreAnswerOne REQUEST: data=%s,meta=%s", request.data, request.meta)
            response.results.append(build_result(request.data))
        return response

    def talkBidirectional(self, request_iterator, context):
        read_headers("TalkBidirectional", context)
        for request in request_iterator:
            logger.info("TalkBidirectional REQUEST: data=%s,meta=%s", request.data, request.meta)
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
