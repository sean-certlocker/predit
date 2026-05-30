package com.predi.app;

import java.util.Date;

public class Market {
    private String id;
    private String title;
    private String status;
    private double currentPool;
    private double userStake;

    public Market(String id, String title, String status) {
        this.id = id;
        this.title = title;
        this.status = status;
    }

    // Getters and Setters
    public String getId() { return id; }
    public String getTitle() { return title; }
    public String getStatus() { return status; }
}
