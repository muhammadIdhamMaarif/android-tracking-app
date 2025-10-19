package com.example.loclogger

import android.Manifest
import android.content.Intent
import android.content.pm.PackageManager
import android.net.Uri
import android.os.Build
import android.os.Bundle
import android.provider.Settings
import androidx.activity.result.contract.ActivityResultContracts
import androidx.appcompat.app.AlertDialog
import androidx.appcompat.app.AppCompatActivity
import androidx.core.app.ActivityCompat
import com.example.loclogger.databinding.ActivityMainBinding

class MainActivity : AppCompatActivity() {

    private lateinit var binding: ActivityMainBinding

    private val requiredPerms = mutableListOf(Manifest.permission.ACCESS_FINE_LOCATION).apply {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.Q) add(Manifest.permission.ACCESS_BACKGROUND_LOCATION)
    }.toTypedArray()

    private val requestPerms =
        registerForActivityResult(ActivityResultContracts.RequestMultiplePermissions()) { perms ->
            val granted = perms.values.all { it }
            if (granted) {
                startLoggingService()
            } else {
                AlertDialog.Builder(this)
                    .setTitle("Permissions required")
                    .setMessage("Location permissions are required for logging. Please enable them in Settings.")
                    .setPositiveButton("Open Settings") { _, _ ->
                        val intent = Intent(Settings.ACTION_APPLICATION_DETAILS_SETTINGS)
                        intent.data = Uri.parse("package:$packageName")
                        startActivity(intent)
                    }
                    .show()
            }
        }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityMainBinding.inflate(layoutInflater)
        setContentView(binding.root)

        binding.btnStart.setOnClickListener { ensurePermissionsAndStart() }
        binding.btnStop.setOnClickListener { stopLoggingService() }
    }

    private fun ensurePermissionsAndStart() {
        val missing = requiredPerms.filter {
            ActivityCompat.checkSelfPermission(this, it) != PackageManager.PERMISSION_GRANTED
        }
        if (missing.isEmpty()) {
            startLoggingService()
        } else {
            requestPerms.launch(requiredPerms)
        }
    }

    private fun startLoggingService() {
        val svcIntent = Intent(this, LocationLoggingService::class.java)
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            startForegroundService(svcIntent)
        } else {
            startService(svcIntent)
        }
        binding.tvStatus.text = "Logging started"
        // prompt user to disable battery optimization for your app
        askDisableBatteryOptimization()
    }

    private fun stopLoggingService() {
        val svcIntent = Intent(this, LocationLoggingService::class.java)
        stopService(svcIntent)
        binding.tvStatus.text = "Logging stopped"
    }

    private fun askDisableBatteryOptimization() {
        // opens battery optimization settings for the app
        val intent = Intent()
        intent.action = Settings.ACTION_IGNORE_BATTERY_OPTIMIZATION_SETTINGS
        startActivity(intent)
        // For direct app-level setting you can use ACTION_REQUEST_IGNORE_BATTERY_OPTIMIZATIONS
        // but it requires manifest and user confirmation.
    }
}
