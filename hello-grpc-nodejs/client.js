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

function main() {
    let client = new proto.LandingService('localhost:9996',
        grpc.credentials.createInsecure());
    client.talk({data: "query=ai,from=0,size=1000,order=x,sort=y", meta: "user=eric"}, function(err, response) {
        console.log('Response:', response);
    });
}

main();