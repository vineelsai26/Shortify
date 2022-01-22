const express = require("express")
const bodyParser = require("body-parser")
const mongodb = require('mongoose')
const router = require('./routers/urlRouter')
const axios = require('axios')
require('dotenv').config()

const app = express()
const mongodbUrl = process.env.MONGODB
const dbRequestUrl = 'http://localhost:5000/generateUrl'

app.set('view engine', 'ejs')
app.use(bodyParser.urlencoded({ extended: false }))
app.use(bodyParser.json())
app.use(express.static('public'))
app.use('/generateUrl', router)

mongodb.connect(mongodbUrl)
const connect = mongodb.connection
connect.on('open', () => {
    console.log('connected')
})

app.get('/', (req, res) => {
    res.render('index', { 'postedUrl': '', 'error': '', 'url': '' })
})

app.post('/', async (req, res) => {
    let postedUrl = req.body.url
    postedUrl = postedUrl.trim()
    postedUrl = postedUrl.toLowerCase()
    postedUrl = postedUrl.replace('https://', '')
    postedUrl = postedUrl.replace('http://', '')

    if (isUrl(postedUrl)) {
        await axios.get(dbRequestUrl)
            .then(res => res.data)
            .then(async (json) => {
                const newURL = await generateUrl(json, postedUrl)
                const url = req.protocol + '://' + req.get('host') + '/' + newURL
                res.render('index', { 'postedUrl': postedUrl, 'error': '', 'url': url })
            })
    } else {
        res.render('index', { 'postedUrl': postedUrl, 'error': 'Invalid URL', 'url': '' })
    }
})

app.get('/:url', (req, res) => {
    const redirectUrl = req.params.url

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
    let regexp = /^(?:(?:https?):\/\/)?(?:(?!(?:10|127)(?:\.\d{1,3}){3})(?!(?:169\.254|192\.168)(?:\.\d{1,3}){2})(?!172\.(?:1[6-9]|2\d|3[0-1])(?:\.\d{1,3}){2})(?:[1-9]\d?|1\d\d|2[01]\d|22[0-3])(?:\.(?:1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.(?:[1-9]\d?|1\d\d|2[0-4]\d|25[0-4]))|(?:(?:[a-z\u00a1-\uffff0-9]-*)*[a-z\u00a1-\uffff0-9]+)(?:\.(?:[a-z\u00a1-\uffff0-9]-*)*[a-z\u00a1-\uffff0-9]+)*(?:\.(?:[a-z\u00a1-\uffff]{2,})))(?::\d{2,5})?(?:\/\S*)?$/
    return regexp.test(str)
}

async function generateUrl(json, postedUrl) {
    let urls = []
    let newUrls = []
    let newURL = ''
    json.forEach((item) => {
        urls.push(item.url)
        newUrls.push(item.newURL)
    })
    if (urls.includes(postedUrl)) {
        json.forEach((item) => {
            if (item.url === postedUrl) {
                newURL = item.newURL
            }
        })
    } else {
        newURL = generateNewUrl(newUrls)

        await axios.post(dbRequestUrl, {
            url: postedUrl,
            newURL: newURL
        }).then(res => {
            console.log(`statusCode: ${res.statusText}`)
        }).catch(error => {
            console.error(error)
        })
    }
    return newURL
}

function generateNewUrl(newUrls) {
    let newURL = makeUrl()
    if (newUrls.includes(newURL)) {
        generateNewUrl(newUrls)
    } else {
        return newURL
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
