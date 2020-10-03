package org.feuyeux.grpc.server;

import io.grpc.ForwardingServerCall.SimpleForwardingServerCall;
import io.grpc.Metadata;
import io.grpc.ServerCall;
import io.grpc.ServerCallHandler;
import io.grpc.ServerInterceptor;
import lombok.extern.slf4j.Slf4j;

@Slf4j
public class HeaderServerInterceptor implements ServerInterceptor {
    static final Metadata.Key<String> k3 = Metadata.Key.of("k3", Metadata.ASCII_STRING_MARSHALLER);

    @Override
    public <ReqT, RespT> ServerCall.Listener<ReqT> interceptCall(
            ServerCall<ReqT, RespT> call,
            final Metadata requestHeaders,
            ServerCallHandler<ReqT, RespT> next) {
        log.info("{} HEADERS:{}", call.getMethodDescriptor().getFullMethodName(), requestHeaders);
        return next.startCall(new SimpleForwardingServerCall<ReqT, RespT>(call) {
            @Override
            public void sendHeaders(Metadata responseHeaders) {
                responseHeaders.put(k3, "v3");
                super.sendHeaders(responseHeaders);
            }
        }, requestHeaders);
    }
}