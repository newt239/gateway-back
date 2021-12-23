import express from 'express';
import verifyToken from 'src/jwt';
var app = express();
var jwt = require('jsonwebtoken');
app.use(express.json());
app.use(express.urlencoded({ extended: true }));

var authRouter = require('./routes/v1/auth/index');

// unauthenticated
app.get('/', function (req: express.Request, res: express.Response) {
    res.json({ status: "OK" });
})

//認証有りAPI
app.get('/protected', verifyToken, function (req: express.Request, res: express.Response) {
    res.json("Protected Contents");
})
app.use('/v1/auth', authRouter);

app.listen(3000, function () {
    console.log("App start on port 3000");
})

module.exports = app;