let grpc = require('grpc')
let sleep = require('sleep')
let messages = require('./common/landing_pb')
let services = require('./common/landing_grpc_pb')

const {createLogger, format, transports} = require('winston')
const {combine, timestamp, printf} = format
const formatter = printf(({level, message, timestamp}) => {
    return `${timestamp} [${level}] ${message}`
})

const logger = createLogger({
    level: 'info',
    format: combine(
        timestamp(),
        formatter
    ),
    transports: [new transports.Console()],
})

function grpcServer() {
    let server = process.env.GRPC_SERVER;
    if (typeof server !== 'undefined' && server !== null) {
        return server
    } else {
        return "localhost"
    }
}

function talk(client, request) {
    logger.info("Talk:" + request)
    client.talk(request, function (err, response) {
        logger.info("Talk:" + response)
        // let result = response.getResultsList()[0]
        // let kvMap = result.getKvMap()
        // logger.info("Talk:" + kvMap.get("data") + response)
    })
}

function talkOneAnswerMore(client, request) {
    logger.info("TalkOneAnswerMore:" + request)
    let call = client.talkOneAnswerMore(request)

    call.on('data', function (response) {
        logger.info("TalkOneAnswerMore:" + response)
        // let result = response.getResultsList()[0]
        // let kvMap = result.getKvMap()
        // logger.info(response)
        // logger.info(kvMap.get("data"))
    })
    call.on('end', function () {
        logger.debug("DONE")
    })
}

function talkMoreAnswerOne(client, requests) {
    let call = client.talkMoreAnswerOne(function (err, response) {
        logger.info("TalkMoreAnswerOne:" + response)
    })
    requests.forEach(request => {
        logger.info("TalkMoreAnswerOne:" + request)
        call.write(request)
    })
    call.end()
}

function talkBidirectional(client, requests) {
    let call = client.talkBidirectional()
    call.on('data', function (response) {
        logger.info("TalkBidirectional:" + response)
    })
    requests.forEach(request => {
        logger.info("TalkBidirectional:" + request)
        sleep.msleep(5)
        call.write(request)
    })
    call.end()
}

function randomId(max) {
    return Math.floor(Math.random() * Math.floor(max)).toString()
}

function main() {
    let address = grpcServer() + ":9996";
    let c = new services.LandingServiceClient(address, grpc.credentials.createInsecure())
    //
    logger.info("Unary RPC")
    let request = new messages.TalkRequest()
    request.setData("0")
    request.setMeta("NODEJS")
    talk(c, request)

    sleep.msleep(50)
    //
    logger.info("Server streaming RPC")
    request = new messages.TalkRequest()
    request.setData("0,1,2")
    request.setMeta("NODEJS")
    talkOneAnswerMore(c, request)

    sleep.msleep(50)
    //
    logger.info("Client streaming RPC")
    let r1 = new messages.TalkRequest()
    r1.setData(randomId(5))
    r1.setMeta("NODEJS")
    let r2 = new messages.TalkRequest()
    r2.setData(randomId(5))
    r2.setMeta("NODEJS")
    let r3 = new messages.TalkRequest()
    r3.setData(randomId(5))
    r3.setMeta("NODEJS")
    let rs = [r1, r2, r3]
    talkMoreAnswerOne(c, rs)

    sleep.msleep(50)
    //
    logger.info("Bidirectional streaming RPC")
    talkBidirectional(c, rs)
}

main()