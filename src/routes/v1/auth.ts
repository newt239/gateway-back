import express from 'express';
import { connectDb } from '@/db';
const router = express.Router();
const jwt = require('jsonwebtoken');

router.post('/login', function (req: express.Request, res: express.Response) {
    const userid: string = req.body.userid;
    const password: string = req.body.password;
    const connection = connectDb(userid, password);
    connection.connect(function (err: any) {
        if (err) {
            res.json({
                status: "error",
                message: "username or password were incorrect.",
                timestamp: Date.now()
            });
        } else {
            const token = jwt.sign({ userid: userid, password: password }, process.env.SIGNATURE);
            res.json({
                status: "success",
                token: token,
                timestamp: Date.now()
            });
        };
    });
});

module.exports = router;