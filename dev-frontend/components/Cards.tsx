import Link from "next/link";
import { ReactElement } from "react";

interface Card {
    title: string,
    content: ReactElement,
    href?: string
}
 
export default function Card({title, content, href}: Card) {
    return (
        <div>
            <div>
                {href 
                    ? <Link href={href}><p>{title}</p></Link> 
                    : <p>{title}</p>}
            </div>
            <div>
                {content}
            </div>
        </div>
    )
}