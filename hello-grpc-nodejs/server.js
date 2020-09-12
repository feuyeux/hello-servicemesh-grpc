let PROTO_PATH = __dirname + '/proto/landing.proto';

let grpc = require('grpc');
let protoLoader = require('@grpc/proto-loader');
let packageDefinition = protoLoader.loadSync(
    PROTO_PATH,
    {keepCase: true,
        longs: String,
        enums: String,
        defaults: true,
        oneofs: true
    });
let proto = grpc.loadPackageDefinition(packageDefinition).org.feuyeux.grpc;

/**
 * Implements the SayHello RPC method.
 */
function talk(call, callback) {
    callback(null, {status: 200});
}
function talkOneAnswerMore(call) {
    //loop
    call.write({status: 200})

    call.end();
}

function talkMoreAnswerOne(call, callback) {

}

function talkBidirectional(call, callback) {

}

/**
 * Starts an RPC server that receives requests for the Greeter service at the
 * sample server port
 */
function main() {
    let server = new grpc.Server();
    server.addService(proto.LandingService.service, {
        talk: talk,
        talkOneAnswerMore: talkOneAnswerMore,
        talkMoreAnswerOne: talkMoreAnswerOne,
        talkBidirectional: talkBidirectional
    });
    server.bind('0.0.0.0:9996', grpc.ServerCredentials.createInsecure());
    server.start();
}

main();