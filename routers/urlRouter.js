const express = require("express")
const router = express.Router()
const UrlModel = require('../models/urlModel')

router.get('/', async (req, res) => {
    try {
        const url = await UrlModel.find()
        res.send(url)
    } catch (err) {
        console.log(err)
    }
})

router.post('/', async (req, res) => {
    const NewUrl = new UrlModel({
        url: req.body.url,
        newURL: req.body.newURL
    })

    try {
        const saveToDB = await NewUrl.save()
        res.json(saveToDB)
    } catch (err) {
        res.send(err)
    }
})

module.exports = router
