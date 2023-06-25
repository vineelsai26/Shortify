import express from 'express'
import bodyParser from 'body-parser'
import mongodb from 'mongoose'
import UrlModel from './models/urlModel.js'
import { nanoid } from 'nanoid'
import 'dotenv/config'
import fs from 'fs'

const app = express()

let mongodbUrl = process.env.MONGODB

if (mongodbUrl === undefined || mongodbUrl?.trim() === '') {
    try {
        const data = fs.readFileSync('/run/secrets/mongodb_url')
        mongodbUrl = data.toString().trim()
    } catch (err) {
        console.log(err)
    }
}

app.set('view engine', 'ejs')
app.use(bodyParser.urlencoded({ extended: false }))
app.use(bodyParser.json())
app.use(express.static('public'))

if (mongodbUrl == undefined || mongodbUrl?.trim() === '') {
    console.log('Please add mongodb url in .env file or mongodb_url file')
    process.exit(1)
}

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
    let protocol = "https"

    if (postedUrl.includes("https://")){
        protocol = "https"
    } else if (postedUrl.includes("http://")) {
        protocol = "http"
    }

    postedUrl = formatUrl(postedUrl)

    if (isUrl(postedUrl)) {
        const newURL = await generateUrl(postedUrl, protocol)
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

function formatUrl(postedUrl) {
    postedUrl = postedUrl.trim()
    postedUrl = postedUrl.replace('https://', '')
    postedUrl = postedUrl.replace('http://', '')
    return postedUrl
}

function isUrl(str) {
    if (str.length < 1000) {
        let regexp = /^(?:(?:https?):\/\/)?(?:(?!(?:10|127)(?:\.\d{1,3}){3})(?!(?:169\.254|192\.168)(?:\.\d{1,3}){2})(?!172\.(?:1[6-9]|2\d|3[0-1])(?:\.\d{1,3}){2})(?:[1-9]\d?|1\d\d|2[01]\d|22[0-3])(?:\.(?:1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.(?:[1-9]\d?|1\d\d|2[0-4]\d|25[0-4]))|(?:(?:[a-z\u00a1-\uffff0-9]-*)*[a-z\u00a1-\uffff0-9]+)(?:\.[a-z\u00a1-\uffff0-9]+)*(?:\.(?:[a-z\u00a1-\uffff]{2,})))(?::\d{2,5})?(?:\/\S*)?$/
        return regexp.test(str)
    } else {
        return false
    }
}

async function generateUrl(redirectUrl, protocol) {
    const url = await UrlModel.findOne({ redirectUrl: redirectUrl })
    if (url && url.url) {
        return url.url
    } else {
        const newURL = await generateNewUrl()

        const NewUrl = new UrlModel({
            url: newURL,
            redirectUrl: redirectUrl,
            protocol: protocol
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
    let newURL = nanoid(6)
    const url = await UrlModel.findOne({ url: newURL })
    if (url && url._id) {
        return await generateNewUrl()
    } else {
        return newURL
    }
}

app.listen(process.env.PORT || 5000)
