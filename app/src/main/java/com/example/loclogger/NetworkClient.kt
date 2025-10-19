package com.example.loclogger

import okhttp3.MediaType.Companion.toMediaTypeOrNull
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.RequestBody.Companion.toRequestBody
import java.util.concurrent.TimeUnit

object NetworkClient {
    private val client = OkHttpClient.Builder()
        .connectTimeout(15, TimeUnit.SECONDS)
        .callTimeout(30, TimeUnit.SECONDS)
        .build()

    // Put your server URL here
    var SERVER_URL = "https://yourserver.example/api/loc"

    fun postJson(json: String): Boolean {
        val body = json.toRequestBody("application/json; charset=utf-8".toMediaTypeOrNull())
        val req = Request.Builder().url(SERVER_URL).post(body).build()
        client.newCall(req).execute().use { resp ->
            return resp.isSuccessful
        }
    }
}
