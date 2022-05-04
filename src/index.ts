import express from 'express';
const fs = require('fs');
const https = require('https');
require('dotenv').config();


const app = express();
app.use(express.json());
app.use(express.urlencoded({ extended: true }));

const cors = require('cors');

const originList = process.env.ORIGIN?.split(',') || [];
const corsConfig = {
    origin: (origin: string, callback: any) => {
        if (originList.indexOf(origin) !== -1) {
            callback(null, true);
        } else {
            callback(new Error('Not allowed by CORS'));
        };
    },
    credentials: true,
    optionsSuccessStatus: 200
};

app.use(cors(corsConfig));

app.get('/', function (req, res) {
    res.json({ status: "OK" });
});

const v1Router = require('./routes/v1/index');
app.use('/v1/', v1Router);

const options = {
    key: fs.readFileSync('/etc/letsencrypt/live/api.sh-fes.com/privkey.pem'),
    cert: fs.readFileSync('/etc/letsencrypt/live/api.sh-fes.com/fullchain.pem'),
}

const server = https.createServer(options, app);
server.listen(443, () => {
    process.setuid && process.setuid('node');
    console.log(`user was replaced to uid: ${process.getuid()} ('node')`);
    console.log('example app listening on port 443!');
});