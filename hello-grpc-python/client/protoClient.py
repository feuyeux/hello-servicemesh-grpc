# encoding: utf-8
import logging
import os
import random
import time

import grpc

from landing_pb2 import landing_pb2
from landing_pb2 import landing_pb2_grpc

logger = logging.getLogger('hello-grpc')
logger.setLevel(logging.INFO)
console = logging.StreamHandler()
console.setLevel(logging.INFO)
formatter = logging.Formatter('%(asctime)s [%(levelname)s] - %(message)s')
console.setFormatter(formatter)
logger.addHandler(console)


def talk(stub):
    request = landing_pb2.TalkRequest(data="0", meta="PYTHON")
    logger.info(request)
    response = stub.talk(request)
    logger.info(response)


def talk_one_answer_more(stub):
    request = landing_pb2.TalkRequest(data="0,1,2", meta="PYTHON")
    logger.info(request)
    responses = stub.talkOneAnswerMore(request)
    for response in responses:
        logger.info(response)


def random_id(end):
    return str(random.randint(0, end))


def generate_request():
    for _ in range(0, 3):
        request = landing_pb2.TalkRequest(data=random_id(5), meta="PYTHON")
        logger.info(request)
        yield request
        time.sleep(random.uniform(0.5, 1.5))


def talk_more_answer_one(stub):
    request_iterator = generate_request()
    response_summary = stub.talkMoreAnswerOne(request_iterator)
    logger.info(response_summary)


def talk_bidirectional(stub):
    request_iterator = generate_request()
    responses = stub.talkBidirectional(request_iterator)
    for response in responses:
        logger.info(response)


def grpc_server():
    server = os.getenv("GRPC_SERVER")
    if server:
        return server
    else:
        return "localhost"


def run():
    with grpc.insecure_channel(grpc_server() + ":9996") as channel:
        stub = landing_pb2_grpc.LandingServiceStub(channel)
        logger.info("Unary RPC")
        talk(stub)
        logger.info("Server streaming RPC")
        talk_one_answer_more(stub)
        logger.info("Client streaming RPC")
        talk_more_answer_one(stub)
        logger.info("Bidirectional streaming RPC")
        talk_bidirectional(stub)


if __name__ == '__main__':
    run()
