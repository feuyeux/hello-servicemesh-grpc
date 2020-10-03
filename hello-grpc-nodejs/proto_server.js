let grpc = require('grpc')
let uuid = require('uuid')
let messages = require('./common/landing_pb')
let services = require('./common/landing_grpc_pb')

let hellos = ["Hello", "Bonjour", "Hola", "こんにちは", "Ciao", "안녕하세요"]

const {createLogger, format, transports} = require('winston')
const {combine, timestamp, printf} = format
const formatter = printf(({level, message, timestamp}) => {
    return `${timestamp} ${level}: ${message}`
})
const logger = createLogger({
    level: 'info',
    format: combine(
        timestamp(),
        formatter
    ),
    transports: [new transports.Console()],
})

function talk(call, callback) {
    printHeaders("Talk", call)
    let response = new messages.TalkResponse()
    response.setStatus(200)
    const talkResult = buildResult(call.request.getData())
    let talkResults = [talkResult]
    response.setResultsList(talkResults)
    logger.info("Talk:" + response)
    callback(null, response)
}

function talkOneAnswerMore(call) {
    printHeaders("TalkOneAnswerMore",call)
    let datas = call.request.getData().split(",")
    for (const data in datas) {
        let response = new messages.TalkResponse()
        response.setStatus(200)
        const talkResult = buildResult(data)
        let talkResults = [talkResult]
        response.setResultsList(talkResults)
        logger.info("TalkOneAnswerMore:" + response)
        call.write(response)
    }
    call.end()
}

function talkMoreAnswerOne(call, callback) {
    printHeaders("TalkMoreAnswerOne",call)
    let talkResults = []
    call.on('data', function (request) {
        talkResults.push(buildResult(request.getData()))
    })
    call.on('end', function () {
        let response = new messages.TalkResponse()
        response.setStatus(200)
        response.setResultsList(talkResults)
        logger.info("TalkMoreAnswerOne:" + response)
        callback(null, response)
    })
}

function talkBidirectional(call) {
    printHeaders("TalkBidirectional",call)
    call.on('data', function (request) {
        let response = new messages.TalkResponse()
        response.setStatus(200)
        let data = request.getData();
        const talkResult = buildResult(data)
        let talkResults = [talkResult]
        response.setResultsList(talkResults)
        logger.info("TalkBidirectional:" + response)
        call.write(response)
    })
    call.on('end', function () {
        call.end()
    })
}

// {"status":200,"results":[{"id":1600402320493411000,"kv":{"data":"Hello","id":"0"}}]}
function buildResult(id) {
    let result = new messages.TalkResult()
    let index = parseInt(id)
    result.setId(Math.round(Date.now() / 1000))
    result.setType(messages.ResultType.OK)
    let kv = result.getKvMap()
    kv.set("id", uuid.v1())
    kv.set("idx", id)
    kv.set("data", hellos[index])
    kv.set("meta", "NODEJS")
    return result
}

/**
 * Starts an RPC server that receives requests for the Greeter service at the
 * sample server port
 */
function main() {
    let server = new grpc.Server()
    server.addService(services.LandingServiceService, {
        talk: talk,
        talkOneAnswerMore: talkOneAnswerMore,
        talkMoreAnswerOne: talkMoreAnswerOne,
        talkBidirectional: talkBidirectional
    })
    server.bind('0.0.0.0:9996', grpc.ServerCredentials.createInsecure())
    server.start()
}

function printHeaders(methodName,call) {
    let headers = call.metadata.getMap();
    for (let key in headers) {
        logger.info(methodName+ " HEADER: " + key + ":" + headers[key])
    }
}

main()