require('dotenv').config()
const { GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET, PORT } = process.env
module.exports = {
    GOOGLE_CLIENT_ID,
    GOOGLE_CLIENT_SECRET,
    PORT
}