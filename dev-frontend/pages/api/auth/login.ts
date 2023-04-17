import type { NextApiRequest, NextApiResponse } from "next";

type Data = {
    uid: string,
    password: string
}

export default function handler(req: NextApiRequest, res: NextApiResponse<Data>) {
    res.status(200).json({
        uid: "12345678901",
        password: "12345678"
    });
}