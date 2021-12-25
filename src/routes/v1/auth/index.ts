import express from 'express';
var jwt = require('jsonwebtoken');
var router = express.Router();

router.post('/login', function (req: express.Request, res: express.Response) {
    var username: string = req.body.username;
    var password: string = req.body.password;
    if (username === "hoge" && password === "password") {
        const token = jwt.sign({ username: username, password: password }, process.env.SIGNATURE);
        res.json({
            status: "success",
            token: token
        });
    } else {
        res.json({
            status: "error",
            message: "auth error"
        });
    }
});

module.exports = router;