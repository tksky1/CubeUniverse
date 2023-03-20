import { ReactElement } from "react"
import MainLayout from "@/components/MainLayout"

export function commonLayout(page: ReactElement) {
    return <MainLayout>{page}</MainLayout>
}