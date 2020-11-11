const mongoose = require('mongoose')


const urlSchema = new mongoose.Schema({

    url: {
        type: String,
        unique: true,
        required: true
    },
    newURL: {
        type: String,
        unique: true,
        required: true
    }
})

module.exports = mongoose.model('Url', urlSchema)
