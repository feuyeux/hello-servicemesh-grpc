import os
import grpc
import sys
sys.path.append(os.path.abspath('.'))
from pb import landing_pb2
from pb import landing_pb2_grpc

def talk(stub):
    request = landing_pb2.TalkRequest(data="query=ai,from=0,size=1000,order=x,sort=y", meta="user=eric")
    print(request)
    response = stub.talk(request)
    print(response)

    result = response.results[0]
    print("First TalkResult: id=%d, type=%s, kv=%s" % (result.id, result.type, result.kv))


def run():
    channel = grpc.insecure_channel('localhost:50061')
    stub = landing_pb2_grpc.LandingServiceStub(channel)
    talk(stub)

if __name__ == '__main__':
    run()
