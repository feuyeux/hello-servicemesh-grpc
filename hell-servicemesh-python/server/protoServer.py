import os,sys
import time
import random
import grpc
from concurrent import futures
sys.path.append(os.path.abspath('.'))
from pb import landing_pb2
from pb import landing_pb2_grpc

_ONE_DAY_IN_SECONDS = 60 * 60 * 24

class LandingServiceServicer(landing_pb2_grpc.LandingServiceServicer):
    def talk(self, request, context):
        return landing_pb2.TalkResponse(
            status=200,
            results=[landing_pb2.TalkResult(
                id=random.randint(0, 9999999),
                type=landing_pb2.SEARCH,
                kv=dict({
                    "request-data": request.data,
                    "request-meta": request.meta,
                    "timestamp": str(time.time())
                })
            )]
        )

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    landing_pb2_grpc.add_LandingServiceServicer_to_server(LandingServiceServicer(), server)
    server.add_insecure_port('[::]:50061')
    server.start()
    try:
        while True:
            time.sleep(_ONE_DAY_IN_SECONDS)
    except KeyboardInterrupt:
        server.stop(0)

if __name__ == '__main__':
    serve()