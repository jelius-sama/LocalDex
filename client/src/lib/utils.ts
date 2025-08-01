import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"
import { useState, useEffect } from "react"

export function cn(...inputs: ClassValue[]) {
    return twMerge(clsx(inputs))
}

export function timeAgo(dateString: string): string {
    const then = new Date(dateString)
    const now = new Date()
    const diff = Math.floor((now.getTime() - then.getTime()) / 1000) // in seconds

    if (diff < 60) return 'just now'
    if (diff < 3600) return `${Math.floor(diff / 60)} minutes ago`
    if (diff < 86400) return `${Math.floor(diff / 3600)} hours ago`
    return `${Math.floor(diff / 86400)} days ago`
}

export function useTimeAgo(dateString?: string): string | null {
    const [value, setValue] = useState<string | null>(() =>
        dateString ? timeAgo(dateString) : null
    )

    useEffect(() => {
        if (!dateString) return

        setValue(timeAgo(dateString))

        const interval = setInterval(() => {
            setValue(timeAgo(dateString))
        }, 60_000)

        return () => clearInterval(interval)
    }, [dateString])

    return value
}

export const formatDate = (dateString: string) => {
    const options: Intl.DateTimeFormatOptions = {
        year: "numeric",
        month: "long",
        day: "numeric",
        hour: "2-digit",
        minute: "2-digit",
    }
    return new Date(dateString).toLocaleDateString(undefined, options)
}

export function formatViews(num: number): string {
    if (num >= 1_000_000_000) {
        return (num / 1_000_000_000).toFixed(1).replace(/\.0$/, '') + 'B';
    }
    if (num >= 1_000_000) {
        return (num / 1_000_000).toFixed(1).replace(/\.0$/, '') + 'M';
    }
    if (num >= 1_000) {
        return (num / 1_000).toFixed(1).replace(/\.0$/, '') + 'K';
    }
    return num.toString();
}
