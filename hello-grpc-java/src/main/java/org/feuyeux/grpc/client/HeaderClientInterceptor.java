package org.feuyeux.grpc.client;

import io.grpc.*;
import lombok.extern.slf4j.Slf4j;

@Slf4j
public class HeaderClientInterceptor implements ClientInterceptor {
    static final Metadata.Key<String> k1 = Metadata.Key.of("k1", Metadata.ASCII_STRING_MARSHALLER);
    static final Metadata.Key<String> k2 = Metadata.Key.of("k2", Metadata.ASCII_STRING_MARSHALLER);

    @Override
    public <ReqT, RespT> ClientCall<ReqT, RespT> interceptCall(MethodDescriptor<ReqT, RespT> method,
                                                               CallOptions callOptions, Channel next) {
        return new ForwardingClientCall.SimpleForwardingClientCall<ReqT, RespT>(next.newCall(method, callOptions)) {
            @Override
            public void start(Listener<RespT> responseListener, Metadata headers) {
                headers.put(k1, "v1");
                headers.put(k2, "v2");
                super.start(new ForwardingClientCallListener.SimpleForwardingClientCallListener<RespT>(responseListener) {
                    @Override
                    public void onHeaders(Metadata headers) {
                        log.info("header received from server:" + headers);
                        super.onHeaders(headers);
                    }
                }, headers);
            }
        };
    }
}