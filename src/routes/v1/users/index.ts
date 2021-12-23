import express from 'express';
var jwt = require('jsonwebtoken');
import verifyToken from '@/jwt';
var router = express.Router();

router.get('/me', verifyToken, function (req: express.Request, res: express.Response) {
    var username: string = req.body.username;
    var password: string = req.body.password;

    if (username === "hoge" && password === "password") {
        const token = jwt.sign({ username: username }, 'my_secret', { expiresIn: '1h' });
        res.json({
            token: token
        });
    } else {
        res.json({
            error: "auth error"
        });
    }
});

module.exports = router;