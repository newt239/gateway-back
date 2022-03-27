import express from 'express';
import verifyToken from '@/jwt';
import { connectDb } from '@/db';
const router = express.Router();

router.post('/create', verifyToken, function (req: express.Request, res: express.Response) {
    const connection = connectDb(res.locals.userid, res.locals.password);
    const userid: string = req.body.userid;
    const password: string = req.body.password;
    let sql = `CREATE USER '${userid}'@'localhost' IDENTIFIED BY '${password}';`;
    sql += `GRANT ON gateway.* TO '${userid}'@'localhost'`;
    sql += `FLUSH PRIVILEGES;`;
});

module.exports = router;