import mongoose from "mongoose"

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

export default mongoose.model('Url', urlSchema)
