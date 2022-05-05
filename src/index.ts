import express from 'express';
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

app.listen(3000);