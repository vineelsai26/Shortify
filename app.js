const serverless = require('serverless-http');
const express = require("express")
const bodyParser = require("body-parser")
const mongodb = require('mongoose')
const UrlModel = require('./models/urlModel')
require('dotenv').config()

const app = express()
const mongodbUrl = process.env.MONGODB

app.set('view engine', 'ejs')
app.use(bodyParser.urlencoded({ extended: false }))
app.use(bodyParser.json())
app.use(express.static('public'))

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
    postedUrl = postedUrl.replace('https://', '')
    postedUrl = postedUrl.replace('http://', '')

    if (isUrl(postedUrl)) {
        const newURL = await generateUrl(postedUrl)
        const url = req.protocol + '://' + req.get('host') + '/' + newURL
        res.render('index', { 'postedUrl': postedUrl, 'error': '', 'url': url })

    } else {
        res.render('index', { 'postedUrl': postedUrl, 'error': 'Invalid URL', 'url': '' })
    }
})

app.get('/:url', (req, res) => {
    const url = req.params.url

    UrlModel.findOne({ url: url }, (err, url) => {
        if (err) {
            console.log(err)
        } else {
            if (url) {
                res.redirect(`${url.protocol}://${url.redirectUrl}`)
            } else {
                res.redirect('/')
            }
        }
    })
})

function isUrl(str) {
    if (str.length < 1000) {
        let regexp = /^(?:(?:https?):\/\/)?(?:(?!(?:10|127)(?:\.\d{1,3}){3})(?!(?:169\.254|192\.168)(?:\.\d{1,3}){2})(?!172\.(?:1[6-9]|2\d|3[0-1])(?:\.\d{1,3}){2})(?:[1-9]\d?|1\d\d|2[01]\d|22[0-3])(?:\.(?:1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.(?:[1-9]\d?|1\d\d|2[0-4]\d|25[0-4]))|(?:(?:[a-z\u00a1-\uffff0-9]-*)*[a-z\u00a1-\uffff0-9]+)(?:\.[a-z\u00a1-\uffff0-9]+)*(?:\.(?:[a-z\u00a1-\uffff]{2,})))(?::\d{2,5})?(?:\/\S*)?$/
        return regexp.test(str)
    } else {
        return false
    }
}

async function generateUrl(redirectUrl) {
    const url = await UrlModel.findOne({ redirectUrl: redirectUrl })
    if (url && url.url) {
        return url.url
    } else {
        const newURL = await generateNewUrl()

        const NewUrl = new UrlModel({
            url: newURL,
            redirectUrl: redirectUrl
        })

        try {
            const saveToDB = await NewUrl.save()
            return saveToDB.url
        } catch (err) {
            console.log(err)
        }
    }
}

async function generateNewUrl() {
    let newURL = makeUrl()
    const url = await UrlModel.findOne({ url: newURL })
    if (url && url._id) {
        return await generateNewUrl()
    } else {
        return newURL
    }
}

function makeUrl() {
    let result = ''
    let characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
    let charactersLength = characters.length
    for (let i = 0; i <= 6; i++) {
        result += characters.charAt(Math.floor(Math.random() * charactersLength))
    }
    return result
}

app.listen(process.env.PORT || 5000)
// module.exports.handler = serverless(app)
