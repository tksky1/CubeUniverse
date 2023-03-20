import Link from "next/link";
import { ReactElement } from "react";

interface Card<T extends boolean> {
    title: string,
    disabled: T,
    content: ReactElement,
    href: T extends true ? string : undefined
}

export default function Card<T extends boolean>({title, disabled, content, href}: Card<T>) {
    return (
        <div>
            <div>
            </div>
        </div>
    )
}