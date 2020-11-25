const express = require("express")
const bodyParser = require("body-parser")
const mongodb = require('mongoose')
const router = require('./routers/urlRouter')
const axios = require('axios')
require('dotenv').config()

const app = express()
const url = process.env.MONGODB

app.use(bodyParser.urlencoded({extended: false}))
app.use(bodyParser.json())
app.use(express.static('public'))
app.use('/generateUrl', router)

mongodb.connect(url, {
    useCreateIndex: true,
    useNewUrlParser: true,
    useUnifiedTopology: true
})
const connect = mongodb.connection
connect.on('open', () => {
    console.log('connected')
})

app.get('/', (req, res) => {
    res.sendFile(__dirname + '/index.html')
})

app.post('/url', async (req, res) => {
    let postedUrl = req.body.url
    const dbRequestUrl = req.protocol + '://' + req.get('host') + '/generateUrl'

    postedUrl = postedUrl.replace('https://', '')
    postedUrl = postedUrl.replace('http://', '')

    if (isUrl(postedUrl)) {
        res.send(await axios.get(dbRequestUrl)
            .then(res => res.data)
            .then((json) => {
                return generateUrl(json, postedUrl, req, res)
            }))
    } else {
        res.send("not an url")
    }
    res.end()
})

app.get('/:url', (req, res) => {
    const redirectUrl = req.params.url
    const dbRequestUrl = req.protocol + '://' + req.get('host') + '/generateUrl'

    axios.get(dbRequestUrl)
        .then(res => res.data)
        .then((json) => {
            json.forEach((item) => {
                if (item.newURL === redirectUrl) {
                    res.redirect('https://' + item.url)
                }
            })
        })
})


function isUrl(str) {
    let regexp = /^(?:(?:https?|ftp):\/\/)?(?:(?!(?:10|127)(?:\.\d{1,3}){3})(?!(?:169\.254|192\.168)(?:\.\d{1,3}){2})(?!172\.(?:1[6-9]|2\d|3[0-1])(?:\.\d{1,3}){2})(?:[1-9]\d?|1\d\d|2[01]\d|22[0-3])(?:\.(?:1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.(?:[1-9]\d?|1\d\d|2[0-4]\d|25[0-4]))|(?:(?:[a-z\u00a1-\uffff0-9]-*)*[a-z\u00a1-\uffff0-9]+)(?:\.(?:[a-z\u00a1-\uffff0-9]-*)*[a-z\u00a1-\uffff0-9]+)*(?:\.(?:[a-z\u00a1-\uffff]{2,})))(?::\d{2,5})?(?:\/\S*)?$/;
    return regexp.test(str);
}

function generateUrl(json, postedUrl, req, res) {
    const dbRequestUrl = req.protocol + '://' + req.get('host') + '/generateUrl'
    let urls = []
    let newUrls = []
    json.forEach((item) => {
        urls.push(item.url)
        newUrls.push(item.newURL)
    })
    if (urls.includes(postedUrl)) {
        json.forEach((item) => {
            if (item.url === postedUrl) {
                console.log(item.newURL)
                res.send(req.protocol + '://' + req.get('host') + '/' + item.newURL)
            }
        })
    } else {
        let newUrl = generateNewUrl(newUrls)

        axios.post(dbRequestUrl, {
            url: postedUrl,
            newURL: newUrl
        }).then(res => {
            console.log(`statusCode: ${res.statusText}`)
        }).catch(error => {
            console.error(error)
        })
        res.send(req.protocol + '://' + req.get('host') + '/redirect/' + newUrl)
    }
}

function generateNewUrl(newUrls) {
    let newUrl = makeUrl()
    if (newUrls.includes(newUrl)) {
        generateNewUrl(newUrls)
    } else {
        return newUrl
    }
}

function makeUrl() {
    let result = ''
    let characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz'
    let charactersLength = characters.length
    for (let i = 0; i < 5; i++) {
        result += characters.charAt(Math.floor(Math.random() * charactersLength))
    }
    return result
}

app.listen(process.env.PORT || 5000)
