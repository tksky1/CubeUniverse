import { createContext } from "react";
import { ReactNode, useState, useEffect } from "react"
import { DataContext } from "./DataContext"
import useWebSocket, { ReadyState } from "react-use-websocket";

export const token = createContext("");

export default function DataProvider({ children }: { children: ReactNode }) {
    let [myData, setMyData] = useState(null);
    let { lastMessage } = useWebSocket("ws://192.168.177.201:30401/api/storage/pvcws");
    useEffect(() => {
        if (lastMessage) {
            let newData = JSON.parse(lastMessage.data) as typeof myData;
            console.log(JSON.parse(lastMessage.data));
            setMyData(newData);
        }
    }, [lastMessage]);
    return (
        <DataContext.Provider value={myData}>
            {children}
        </DataContext.Provider>
    )
}