import express from 'express';
import verifyToken from '@/jwt';
import { connectDb } from '@/db';
const router = express.Router();

router.post('/create', verifyToken, function (req: express.Request, res: express.Response) {
    const connection = connectDb(res.locals.userId, res.locals.password);
    const userId: string = req.body.userId;
    const password: string = req.body.password;
    let sql = `CREATE USER '${userId}'@'localhost' IDENTIFIED BY '${password}';`;
    sql += `GRANT ON gateway.* TO '${userId}'@'localhost'`;
    sql += `FLUSH PRIVILEGES;`;
});

module.exports = router;