import express from 'express';
require('dotenv').config();

const app = express();
app.use(express.json());
app.use(express.urlencoded({ extended: true }));

const cors = require('cors');
app.use(cors({
    origin: process.env.ORIGIN,
    credentials: true,
    optionsSuccessStatus: 200
}));

app.get('/', function (req: express.Request, res: express.Response) {
    res.json({ status: "OK" });
});

const v1Router = require('./routes/v1/index');
app.use('/v1/', v1Router);

app.listen(3000);