import express from 'express';
var app = express();
require('dotenv').config();
const cors = require('cors');
app.use(express.json());
app.use(express.urlencoded({ extended: true }));
app.use(cors({
    origin: process.env.ORIGIN,
    credentials: true,
    optionsSuccessStatus: 200
}));

app.get('/', function (req: express.Request, res: express.Response) {
    res.json({ status: "OK" });
})

const authRouter = require('./routes/v1/auth/index');
app.use('/v1/auth/', authRouter);

const usersRouter = require('./routes/v1/users/index');
app.use('/v1/users/', usersRouter);

app.listen(3000, function () {
    console.log("App start on port 3000");
})

module.exports = app;