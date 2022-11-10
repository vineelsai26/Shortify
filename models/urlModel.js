const mongoose = require('mongoose')

const urlSchema = new mongoose.Schema({
    url: {
        type: String,
        unique: true,
        required: true
    },
    redirectUrl: {
        type: String,
        unique: true,
        required: true
    },
    createdAt: {
        type: Date,
        default: Date.now
    },
    protocol: {
        type: String,
        required: true,
        default: 'https'
    }
})

module.exports = mongoose.model('Url', urlSchema)
