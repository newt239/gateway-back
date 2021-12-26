import express from 'express';
var jwt = require('jsonwebtoken');

function verifyToken(req: express.Request, res: express.Response, next: express.NextFunction) {
    const authHeader = req.headers["authorization"];
    if (authHeader !== undefined) {
        if (authHeader.split(" ")[0] === "Bearer") {
            jwt.verify(authHeader.split(" ")[1], process.env.SIGNATURE, (err: any, payload: any) => {
                if (err) {
                    return res.sendStatus(403);
                } else {
                    res.locals.userid = payload.userid;
                    res.locals.password = payload.password;
                    next();
                }
            });
        } else {
            res.json({ error: "header format error" });
        }
    } else {
        res.json({ error: "header error" });
    }
}
export default verifyToken;