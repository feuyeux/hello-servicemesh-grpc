package org.feuyeux.grpc.server;

import io.grpc.stub.StreamObserver;
import lombok.extern.slf4j.Slf4j;
import org.feuyeux.grpc.proto.*;

import java.util.*;

@Slf4j
public class LandingServiceImpl extends LandingServiceGrpc.LandingServiceImplBase {
    private final List<String> HELLO_LIST = Arrays.asList("Hello", "Bonjour", "Hola", "こんにちは", "Ciao", "안녕하세요");

    @Override
    public void talk(TalkRequest request, StreamObserver<TalkResponse> responseObserver) {
        log.info("REQUEST:{}", request.toString());
        String data = request.getData();
        TalkResponse response = TalkResponse.newBuilder()
                .setStatus(200)
                .addResults(buildResult(data)).build();
        log.info("RESPONSE:{}", response.toString());
        responseObserver.onNext(response);
        responseObserver.onCompleted();
    }

    @Override
    public void talkOneAnswerMore(TalkRequest request, StreamObserver<TalkResponse> responseObserver) {
        List<TalkResponse> talkResponses = new ArrayList<>();
        String data0 = request.getData();
        String[] datas = data0.split(",");
        for (String data : datas) {
            TalkResponse response = TalkResponse.newBuilder()
                    .setStatus(200)
                    .addResults(buildResult(data)).build();
            talkResponses.add(response);
        }
        talkResponses.forEach(responseObserver::onNext);
        responseObserver.onCompleted();
    }

    @Override
    public StreamObserver<TalkRequest> talkMoreAnswerOne(StreamObserver<TalkResponse> responseObserver) {
        return new StreamObserver<TalkRequest>() {
            final List<TalkRequest> talkRequests = new ArrayList<>();

            @Override
            public void onNext(TalkRequest request) {
                talkRequests.add(request);
            }

            @Override
            public void onError(Throwable t) {
                log.error("talkBidirectional onError");
            }

            @Override
            public void onCompleted() {
                responseObserver.onNext(buildResponse(talkRequests));
                responseObserver.onCompleted();
            }
        };
    }

    @Override
    public StreamObserver<TalkRequest> talkBidirectional(StreamObserver<TalkResponse> responseObserver) {
        return new StreamObserver<TalkRequest>() {
            @Override
            public void onNext(TalkRequest request) {
                responseObserver.onNext(TalkResponse.newBuilder()
                        .setStatus(200)
                        .addResults(buildResult(request.getData())).build());
            }

            @Override
            public void onError(Throwable t) {
                log.error("talkBidirectional onError");
            }

            @Override
            public void onCompleted() {
                responseObserver.onCompleted();
            }
        };
    }

    private TalkResult buildResult(String id) {
        int index;
        try {
            index = Integer.parseInt(id);
        } catch (NumberFormatException ignored) {
            index = 0;
        }
        String data;
        if (index > 5) {
            data = "你好";
        } else {
            data = HELLO_LIST.get(index);
        }
        Map<String, String> kv = new HashMap<>();
        kv.put("data", data);
        kv.put("id", id);
        TalkResult result = TalkResult.newBuilder()
                .setId(System.nanoTime())
                .setType(ResultType.OK)
                .putAllKv(kv)
                .build();
        return result;
    }

    private TalkResponse buildResponse(List<TalkRequest> talkRequests) {
        TalkResponse.Builder response = TalkResponse.newBuilder();
        response.setStatus(200);
        talkRequests.forEach(request -> {
            response.addResults(buildResult(request.getData()));
        });
        return response.build();
    }
}
