import express from 'express';
import verifyToken from '@/jwt';
import { connectDb } from '@/db'
var router = express.Router();

router.get('/me', verifyToken, function (req: express.Request, res: express.Response) {
    const connection = connectDb(res.locals.userid, res.locals.password);
    return res.json({ userid: res.locals.userid });
});

module.exports = router;