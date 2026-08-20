export type Project = {id:string;code:string;name:string;year:number;status:string}
export type Claim = {id:string;summary:string;amountCents:number;status:string}
export const apiBase='/api/v1'
export function formatCents(v:number):string{return (v/100).toFixed(2)}
export function createAppShell(){return {name:'公益补助工作台',routes:['rules','claims','preview','review','ledger','compare','export']}}
export function piniaStore<T>(initial:T){let state=initial;return {get:()=>state,set:(next:T)=>{state=next}}}
